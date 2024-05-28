<script lang="ts">
  import ChatBubble from "./lib/ChatBubble.svelte";
  import ChatContainer from "./lib/ChatContainer.svelte";
  import ChatHeader from "./lib/ChatHeader.svelte";
  import ChatInput from "./lib/ChatInput.svelte";
  import Container from "./lib/Container.svelte";
  import KeyInputs from "./lib/KeyInputs.svelte";
  import { onMount } from "svelte";
  import { decryptMessage, deriveSharedSecret, encryptMessage, generateKeyPair, Point, bigIntToHex } from "./utils/ecc";
  import wasm from "./wasm/main.go";
  import type { Signature, SchnorrKeys } from "./wasm/main.go";


  type Message = {
    sender: string
    message: string
    verified: boolean
  }

  let messages: Message[] = []
  let id: string

  let privKeyECDH: bigint | null
  let pubKeyECDH: Point | null
  let sharedKeyECDH: Point | null
  let sharedKeyHash: string | null

  let privKeyECC: bigint | null
  let pubKeyECC: Point | null

  let socket: WebSocket
  let isConnected = false

  let status: string
  let error: string

  type Schnorr = {
    p: string
    q: string
    gen: string
  }

  type PublicKey = {
    x: string
    y: string
  }

  let schnorr: Schnorr | null
  let schnorrKeys: SchnorrKeys | null

  const setupSchnorr = async () => {
    const sch: Schnorr = await fetch("http://localhost:8080/schnorr", { method: "GET" })
      .then(response => response.json())
      .then(data => data)
      .catch(error => console.log("error", error))
    // console.log(sch)

    const keys: SchnorrKeys = await wasm.keys(sch.p, sch.q, sch.gen);

    schnorr = sch
    schnorrKeys = keys
    console.log("Schnorr: ", schnorrKeys)
  }

  const signMessage = async (message: string) : Promise<Signature | null> => {
    console.log("Signing message: ", message, message.length)
    if (!schnorr || !localSigningKey) {
      console.error("No schnorr :'((")
      return null
    }
    let s = await wasm.sign(schnorr.p, schnorr.q, schnorr.gen, localSigningKey!, message);
    console.log("Signature: ", s);
    return s
  }

  const verifyMessage = async (message: string, signature: Signature) : Promise<boolean> => {
    console.log("Verifying message: ", message, message.length)
    if (!schnorr || !remotePublicKey) {
      console.error("No schnorr :'((")
      return false
    }
    let v = await wasm.verify(schnorr.p, schnorr.q, schnorr.gen, remotePublicKey!, signature.sign, signature.hash, message);
    console.log("Verified: ", v);
    return v
  }

  let sign: boolean = false
  function doSign(e: Event) {
    sign = (e.target as HTMLInputElement).checked
    console.log("Signing? ", sign)
  }

  let localSigningKey: string | null
  function setSignKey(e: Event) {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (!file) return
    const reader = new FileReader()
    reader.onload = () => {
      localSigningKey = reader.result as string
      console.log("localSigningKey", localSigningKey)
    }
    reader.readAsText(file)
  }

  let remotePublicKey: string | null
  function setVerifyKey(e: Event) {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (!file) return
    const reader = new FileReader()
    reader.onload = () => {
      remotePublicKey = reader.result as string
      console.log("remotePublicKey", remotePublicKey)
    }
    reader.readAsText(file)
  }

  function onGenerateSign() {
    const privKeyBlob = new Blob([schnorrKeys!.private], {type: "text/plain"})
    const pubKeyBlob = new Blob([schnorrKeys!.public], {type: "text/plain"})
    const privKeyURL = URL.createObjectURL(privKeyBlob)
    const pubKeyURL = URL.createObjectURL(pubKeyBlob)
    const privKeyLink = document.createElement("a")
    const pubKeyLink = document.createElement("a")
    privKeyLink.href = privKeyURL
    privKeyLink.download = ".schprv"
    pubKeyLink.href = pubKeyURL
    pubKeyLink.download = ".schpub"
    privKeyLink.click()
    pubKeyLink.click()
    URL.revokeObjectURL(privKeyURL)
    URL.revokeObjectURL(pubKeyURL)
    privKeyLink.remove()
    pubKeyLink.remove()

  }

  // let sharedKey: string | null
  const doECDH = async () => {
    const keyPair = generateKeyPair()
    privKeyECDH = keyPair[0]
    pubKeyECDH = keyPair[1]
    console.log("privKeyECDH", privKeyECDH)
    console.log("id", id)
    let opts = {
      method: "PUT",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        port: id,
        public_key: pointToJSON(pubKeyECDH)
      })
    }

    console.log("opts", opts.body)
    const data: PublicKey = await fetch("http://localhost:8080/key", opts)
      .then(response => response.json())
      .then(data => data)
      .catch(error => console.log("error", error))

    const pubKey = new Point(BigInt(data.x), BigInt(data.y))
    sharedKeyECDH = deriveSharedSecret(privKeyECDH as bigint, pubKey)
    wasm.hash(bigIntToHex(sharedKeyECDH.x)).then((hash) => {
      sharedKeyHash = hash
      console.log("sharedKeyHash", sharedKeyHash)
    })
  }

  // Encrypt and Decrypt messages
  const cipher = async (message: string, isEncrypt: boolean) : Promise<string> => {
    // console.log(sharedKeyECDH)
    let key = sharedKeyHash ?? ""

    if (!key) {
      // TODO: Toast errors
      console.log("No key")
      return Promise.reject("No key")
    }
    // console.log(key, message)

    if (isEncrypt) {
      return await wasm.encrypt(key, message);
    } else {
      return await wasm.decrypt(key, message);
    }
  }

  // Connect to WebSocket server
  const connectWS = () => {
    socket = new WebSocket("ws://localhost:8080/chat?id=" + id)
    socket.addEventListener("open", ()=> {
      console.log("Opened")
      isConnected = true
    })
    socket.addEventListener("message", (event) => {
      cipher(event.data, false)
        .then((text) => {
          console.log("Decrypted: ", text)
          const payload = JSON.parse(text)
          if (payload.sender === id) return
          const decrypted = decryptMessage(privKeyECC as bigint, JSONToPoints(payload.message))
          if (payload.sign && payload.hash && remotePublicKey) {
            verifyMessage(decrypted, {sign: payload.sign, hash: payload.hash})
              .then((verified) => {
                const message = {
                  sender: payload.sender,
                  message: decrypted,
                  verified: verified,
                }
                messages = [message, ...messages]
                if (!verified) {
                  console.log("Message not verified, chance to be tampered with.")
                }
              })
          } else {
            const message = {
              sender: payload.sender,
              message: decrypted,
              verified: false,
            }
            messages = [message, ...messages]
          }
        })
    })
    socket.addEventListener("close", () => {
      console.log("Closed")
      isConnected = false
      // While connection is closed, try to reconnect every 3 seconds
      setTimeout(() => {
        connectWS()
      }, 3000)
    })
    socket.addEventListener("error", (error) => {
      console.log("Error ", error)
    })
  }

  // Convert point to JSON string
  function pointToJSON(point: Point): string {
    return JSON.stringify({x: point.x.toString(), y: point.y.toString()})
  }

  // Convert points to JSON string
  function pointsToJSON(point: [Point, Point]): string {
    return JSON.stringify([{x: point[0].x.toString(), y: point[0].y.toString()}, {x: point[1].x.toString(), y: point[1].y.toString()}])
  }

  // Convert JSON string to points
  function JSONToPoints(json: string): [Point, Point] {
    const points = JSON.parse(json)
    return [new Point(BigInt(points[0].x), BigInt(points[0].y)), new Point(BigInt(points[1].x), BigInt(points[1].y))]
  }

  // Convert JSON string to point
  function JSONToPoint(json: string): Point {
    const point = JSON.parse(json)
    return new Point(BigInt(point.x), BigInt(point.y))
  }

  // Generate key pairs
  function generate() {
    const [priv, pub] = generateKeyPair()
    privKeyECC = priv
    // pubKeyECC = pub

    // download keys
    const privKey = priv.toString()
    const pubKey = pointToJSON(pub)
    const privKeyBlob = new Blob([privKey], {type: "text/plain"})
    const pubKeyBlob = new Blob([pubKey], {type: "text/plain"})
    const privKeyURL = URL.createObjectURL(privKeyBlob)
    const pubKeyURL = URL.createObjectURL(pubKeyBlob)
    const privKeyLink = document.createElement("a")
    const pubKeyLink = document.createElement("a")
    privKeyLink.href = privKeyURL
    privKeyLink.download = ".ecprv"
    pubKeyLink.href = pubKeyURL
    pubKeyLink.download = ".ecpub"
    privKeyLink.click()
    pubKeyLink.click()
    URL.revokeObjectURL(privKeyURL)
    URL.revokeObjectURL(pubKeyURL)
    privKeyLink.remove()
    pubKeyLink.remove()

    status = "Keys generated and downloaded successfully. Please give your public key to your partner."
  }

  function setPrivKeyECC(e: Event) {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (!file) return
    const reader = new FileReader()
    reader.onload = () => {
      privKeyECC = BigInt(reader.result as string)
      console.log("privKeyECC", privKeyECC)
    }
    reader.readAsText(file)
  }

  function setPubKeyECC(e: Event) {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (!file) return
    const reader = new FileReader()
    reader.onload = () => {
      pubKeyECC = JSONToPoint(reader.result as string)
      console.log("pubKeyECC", pubKeyECC)
    }
    reader.readAsText(file)
  }

  onMount(() => {
    // // Try to get private key from local storage
    // privKey = localStorage.getItem("privKey") ? BigInt(localStorage.getItem("privKey") as string) : null
    // pubKey = localStorage.getItem("pubKey") ? JSON.parse(localStorage.getItem("pubKey") as string) : null
    // if (!privKey || !pubKey) {
    //   // If private key is not found, generate a new one
    //   privKey = generatePrivateKey()
    //   pubKey = generatePublicKey(privKey)
    //   localStorage.setItem("privKey", privKey.toString())
    //   // localStorage.setItem("pubKey", JSON.stringify(pubKey))
    // } 
    // console.log("privKey", privKey)
    // console.log("pubKey", pubKey)

    const url = window.location.href
    id = url.split(":")[2].split("/")[0]

    setupSchnorr().then(() => {
      localSigningKey = schnorrKeys!.private
      console.log("Local Signing Key: ", localSigningKey)
    })

    doECDH()
      .then(() => connectWS())


    return () => {
      socket.close()
    }
  })

  const testing = (e: Event) => {
    if (!e.target) return
    console.log((e.target as HTMLInputElement).value)

    const privKeyyy = BigInt((e.target as HTMLInputElement).value)
    console.log("shared Key 2", deriveSharedSecret(privKeyECDH as bigint, pubKeyECDH as Point))
  }

  const onSend = (message: string) => {
    if (!pubKeyECC) {
      error = "Please provide keys."
      return
    }

    if (sign && localSigningKey) {
      signMessage(message)
        .then((signature) => {
          const payload = {
            sender: id,
            message: pointsToJSON(encryptMessage(pubKeyECC!, message)),
            sign: signature?.sign,
            hash: signature?.hash
          }
          const payloadString = JSON.stringify(payload)
          console.log("Sending ", payloadString)

          cipher(payloadString, true)
            .then((text) => {
              console.log("Encrypted: ", text)
              socket.send(text)
              messages = [{sender: id, message: message, verified: true}, ...messages]
            })
        })
    } else {
      const payload = {
        sender: id,
        message: pointsToJSON(encryptMessage(pubKeyECC, message))
      }
      const payloadString = JSON.stringify(payload)
      console.log("Sending ", payloadString)

      cipher(payloadString, true)
        .then((text) => {
          console.log("Encrypted: ", text)
          socket.send(text)
          messages = [{sender: id, message: message, verified: false}, ...messages]
        })
      
    }
  }
</script>

<main class="bg-neutral-100 h-screen">
  <div class="flex flex-col md:flex-row w-full h-screen">
    <KeyInputs doSign={doSign} setSignKey={setSignKey} setVerifyKey={setVerifyKey} onGenerateSign={onGenerateSign} onGenerate={generate} setPrivKeyECC={setPrivKeyECC} setPubKeyECC={setPubKeyECC} status={status} />
    <Container>
      <ChatHeader sender={id} isConnected={isConnected} />
      <ChatContainer>
        <!-- loop with index -->
        {#each messages as message, i}
          <ChatBubble isSelf={message.sender === id} message={message.message} isVerified={message.verified} />
        {/each}
      </ChatContainer>
      <ChatInput onSend={onSend} />
    </Container>
    <div class="w-0 lg:w-1/3">
      {#if error}
        <p class="text-red">{error}</p>
      {/if}
      <input type="text" on:change={testing} />
    </div>
  </div>
</main>
