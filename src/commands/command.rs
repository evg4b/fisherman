use crate::context::Context;
use anyhow::Result;

pub trait CliCommand {
    fn exec(&self, context: &mut impl Context) -> Result<()>;
}
