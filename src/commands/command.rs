use anyhow::Result;
use fisherman_core::Context;

pub trait CliCommand {
    async fn exec(&self, context: &mut impl Context) -> Result<()>;
}
