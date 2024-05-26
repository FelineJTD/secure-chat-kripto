type Signature = {
    sign: string;
    hash: string;
}

type SchnorrKeys = {
    private: string;
    public: string;
}

const __default: {
    keys(p: string, q: string, g: string): Promise<SchnorrKeys>;
    sign(p: string, q: string, g: string, privkey: string, message: string): Promise<Signature>;
    verify(p: string, q: string, g: string, pubkey: string, sign: string, hash: string, message: string): Promise<boolean>;
    encrypt(key: string, plaintext: string): Promise<string>;
    decrypt(key: string, ciphertext: string): Promise<string>;
    hash(hexString: string): Promise<string>;
}

export default __default;
export { Signature, SchnorrKeys };