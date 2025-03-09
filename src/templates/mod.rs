mod placeholders;
mod errors;
mod hashmap;

pub use crate::templates::errors::TemplateError;
pub use crate::templates::hashmap::replace_in_hashmap;
pub use crate::templates::placeholders::replace_in_string;
