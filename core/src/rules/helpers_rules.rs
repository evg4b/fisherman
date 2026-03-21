#[macro_export]
macro_rules! extract_vars {
    ($self:expr, $ctx:expr) => {
        $ctx.variables($self.extract.clone().unwrap_or_default().as_slice())
    };
}
