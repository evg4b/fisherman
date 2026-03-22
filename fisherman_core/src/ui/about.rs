use clap::builder::{IntoResettable, Resettable, StyledStr};

use crate::ui::logo::logo;

pub struct About {}

impl IntoResettable<StyledStr> for About {
    fn into_resettable(self) -> clap::builder::Resettable<StyledStr> {
        Resettable::Value(StyledStr::from(logo()))
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use clap::builder::Resettable;

    #[test]
    fn about_into_resettable_returns_logo() {
        let about = About {};
        let res = about.into_resettable();

        match res {
            Resettable::Value(value) => {
                assert_eq!(value.to_string(), logo());
            }
            other => {
                panic!("Expected Resettable::Value, got {:?}", other);
            }
        }
    }
}