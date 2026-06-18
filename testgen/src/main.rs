use msgpack_tagged::{EncodingStrategy, MsgpackTagged, Serializer, TagRegistry};
use serde::Serialize;

mod acir_field;
mod black_box_func;
mod expression;
mod opcodes;
mod shared;

fn main() {
    let target_dir = "../go/binaries/";

    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::TRACE)
        .with_target(false)
        .with_writer(std::io::stdout)
        .init();

    // check if the directory exists
    if !std::path::Path::new(target_dir).exists() {
        std::fs::create_dir_all(target_dir).expect("Failed to create directory");
    }

    black_box_func::generate_tests(target_dir);
    shared::generate_tests(target_dir);
    acir_field::generate_tests(target_dir);
    expression::generate_tests(target_dir);
    opcodes::generate_tests(target_dir);
}

pub(crate) fn encode<T>(value: &T) -> Vec<u8>
where
    T: ?Sized + Serialize + MsgpackTagged,
{
    let registry = TagRegistry::from_type::<T>();
    let mut buf = Vec::new();
    let mut serializer = Serializer::new(&mut buf, &registry)
        .with_default_strategy(EncodingStrategy::Array)
        .with_strategy_for_name("Program", EncodingStrategy::Tagged)
        .with_strategy_for_name("Circuit", EncodingStrategy::Tagged)
        .with_strategy_for_name("BrilligBytecode", EncodingStrategy::Tagged);
    value
        .serialize(&mut serializer)
        .expect("msgpack-tagged serialize");
    buf
}
