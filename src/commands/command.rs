use anyhow::Result;
use core::Context;

pub trait CliCommand {
    fn exec(&self, context: &mut impl Context) -> Result<()>;
}
