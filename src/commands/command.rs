use anyhow::Result;
use fisherman_core::Context;

pub trait CliCommand {
    fn exec(&self, context: &mut impl Context) -> Result<()>;
}
