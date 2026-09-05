mod errors;
mod hashmap;
mod template;
mod template_macro;
mod vec;

pub use crate::errors::TemplateError;
pub use crate::hashmap::replace_in_hashmap;
pub use crate::template::TemplateString;
pub use crate::vec::replace_in_vec;
