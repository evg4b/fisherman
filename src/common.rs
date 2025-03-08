use std::fmt::{Display, Formatter};

/// Boxed error to be used in the configuration module
#[macro_export]
macro_rules! err {
    ($e:expr) => {
        return Err(Box::new($e))
    };
}

pub type BError = Box<dyn std::error::Error>;
pub type R<T = ()> = Result<T, BError>;

#[derive(Debug)]
pub struct CommonError {
    message: String,
}

impl Display for CommonError {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.message)
    }
}

impl std::error::Error for CommonError {}

impl CommonError {
    pub fn new<T: Into<String>>(message: T) -> Self {
        Self { message: message.into() }
    }
}