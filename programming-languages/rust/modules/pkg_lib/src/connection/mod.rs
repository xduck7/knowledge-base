pub mod status;

use status::Status;

pub(crate) struct Conn {
    max_conn: i8,
    driver: String,
    active: Status,
}

impl Conn {
    pub(crate) fn init(max_conn: i8) -> Conn {
        Conn {
            max_conn,
            driver: "Rust".to_string(),
            active: Status::Connected,
        }
    }

    fn shutdown(&mut self) {
        self.max_conn = 0;
        self.active = Status::Disconnected;
    }

    fn exec(&mut self, cmd: &str) {
        self.active = Status::Busy;
        println!("{}: success execution!", cmd);
        self.active = Status::Connected;
    }

    fn get_max_conn(&self) -> i8 {
        self.max_conn
    }
}
