use anyhow::Result;
use crate::Context;

pub trait CliCommand {
    fn exec(&self, context: &mut impl Context) -> Result<()>;
}
