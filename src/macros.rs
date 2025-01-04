/// Boxed error to be used in the configuration module
#[macro_export]
macro_rules! err {
    ($e:expr) => {
        return Err(Box::new($e))
    };
}