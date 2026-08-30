use crate::Context;
use anyhow::Result;

pub trait CliCommand {
    async fn exec(&self, context: &mut impl Context) -> Result<()>;
}
