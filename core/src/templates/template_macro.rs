#[macro_export]
macro_rules! t {
    ($template:expr) => {
        $crate::TemplateString::from($template.to_string())
    };
}

#[macro_export]
macro_rules! tmpl {
    ($template:expr) => {
        $crate::TemplateString::from($template.to_string())
    };
}
