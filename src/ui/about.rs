use clap::builder::{IntoResettable, Resettable, StyledStr};

use crate::ui::logo::logo;

pub struct About {}

impl IntoResettable<StyledStr> for About {
    fn into_resettable(self) -> clap::builder::Resettable<StyledStr> {
        Resettable::Value(StyledStr::from(logo()))
    }
}