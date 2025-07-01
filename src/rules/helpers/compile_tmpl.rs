use crate::context::Context;
use crate::templates::TemplateString;
use anyhow::Result;

pub fn compile_tmpl(
    ctx: &dyn Context,
    string: &TemplateString,
    additional: &[String],
) -> Result<String> {
    let variables = ctx.variables(additional)?;
    Ok(string.to_string(&variables)?)
}
