#[cfg(test)]
mod tests {
    use std::fs::File;

    use ark_ec::CurveGroup;

    use crate::{
        proof::GnarkProof,
        verifier::{g2_from_bytes, g2_to_bytes, GnarkVerifier},
        vk::parse_vk,
        witness::GnarkWitness,
    };

    #[test]
    fn test_verification_no_commitment() {
        const NR_INPUTS: usize = 1;

        let vk_file = File::open("src/test_files/sum_a_b.vk").expect("unable to open vk file");
        let vk = parse_vk(vk_file).expect("Unable to parse vk");
        let mut verifier = GnarkVerifier::<'_, NR_INPUTS>::new(&vk);

        let proof_file =
            File::open("src/test_files/sum_a_b.proof").expect("unable to open proof file");
        let proof = GnarkProof::parse(proof_file).expect("Unable to parse proof");

        let pw_file = File::open("src/test_files/sum_a_b.pw").expect("unable to open pw file");

        let public_inputs =
            GnarkWitness::<NR_INPUTS>::parse(pw_file).expect("Unable to parse public witness");

        let res = verifier.verify(proof, public_inputs);

        assert!(res.is_ok())
    }
    #[test]
    fn test_verification_with_commitment() {
        const NR_INPUTS: usize = 0;

        let vk_file = File::open("src/test_files/keccak_f1600.vk").expect("unable to open vk file");
        let vk = parse_vk(vk_file).expect("Unable to parse vk");
        let mut verifier = GnarkVerifier::<'_, NR_INPUTS>::new(&vk);

        let proof_file =
            File::open("src/test_files/keccak_f1600.proof").expect("unable to open proof file");
        let proof = GnarkProof::parse(proof_file).expect("Unable to parse proof");

        let pw_file = File::open("src/test_files/keccak_f1600.pw").expect("unable to open pw file");

        let public_inputs =
            GnarkWitness::<NR_INPUTS>::parse(pw_file).expect("Unable to parse public witness");

        let res = verifier.verify(proof, public_inputs);

        assert!(res.is_ok())
    }

    #[test]
    fn test_verification_with_commitment_and_pw() {
        const NR_INPUTS: usize = 2;

        let vk_file = File::open("src/test_files/xor.vk").expect("unable to open vk file");
        let vk = parse_vk(vk_file).expect("Unable to parse vk");
        let mut verifier = GnarkVerifier::<'_, NR_INPUTS>::new(&vk);

        let proof_file = File::open("src/test_files/xor.proof").expect("unable to open proof file");
        let proof = GnarkProof::parse(proof_file).expect("Unable to parse proof");

        let pw_file = File::open("src/test_files/xor.pw").expect("unable to open pw file");

        let public_inputs =
            GnarkWitness::<NR_INPUTS>::parse(pw_file).expect("Unable to parse public witness");

        let res = verifier.verify(proof, public_inputs);

        assert!(res.is_ok())
    }

    #[test]
    fn test_g2_bytes_serde() {
        use ark_bn254::G2Projective;
        use ark_ff::UniformRand;
        use ark_std::rand::thread_rng;

        let mut rng = thread_rng();
        let p = G2Projective::rand(&mut rng).into_affine();
        let bytes = g2_to_bytes(&p);
        let q = g2_from_bytes(&bytes);

        assert_eq!(p, q, "g2 serde failed");
    }
}
