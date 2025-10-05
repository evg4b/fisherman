#[macro_export]
macro_rules! t {
    ($template:expr) => {
        $crate::templates::TemplateString::from($template.to_string())
    };
}
