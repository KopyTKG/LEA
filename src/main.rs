use rand::RngCore;
use rand::rngs::OsRng;

struct LEA {
    plain_text: String,
    key: [usize; 4],
    key_schedule: [[usize; 6]; 24],
    constants: [usize; 4],
    // cipher_text: String,
    // decrypted_text: String,
}

impl LEA {
    fn generate_key_schedule(&self) -> [[usize; 6]; 24] {
        let mut key_schedule: [[usize; 6]; 24] = [[0; 6]; 24];
        let mut t: [usize; 4] = self.key;
    
        for i in 0..24 {
            // Update t for each round
            let mut t_temp = t;
            t_temp[0] = (((t_temp[0] + (self.constants[i % 4] << (i + 0))) << 1) & ((1 << 32) - 1));
            t_temp[1] = (((t_temp[1] + (self.constants[i % 4] << (i + 1))) << 3) & ((1 << 32) - 1));
            t_temp[2] = (((t_temp[2] + (self.constants[i % 4] << (i + 2))) << 6) & ((1 << 32) - 1));
            t_temp[3] = (((t_temp[3] + 1 + (self.constants[i % 4] << (i + 3))) << 11) & ((1 << 32) - 1));
    
            // Assign values to key_schedule
            key_schedule[i][0] = t_temp[0];
            key_schedule[i][1] = t_temp[1];
            key_schedule[i][2] = t_temp[2];
            key_schedule[i][3] = t_temp[1]; // Using t_temp[1] here instead of t[1]
            key_schedule[i][4] = t_temp[3];
            key_schedule[i][5] = t_temp[1]; // Using t_temp[1] here instead of t[1]
    
            // Update t for the next round
            t = t_temp;
        }
    
        return key_schedule;
    }
        
    fn explode_plaintext(&self) -> [usize; 4] {
        let mut segments: [usize; 4] = [0; 4];
        
        // Convert the plaintext characters into their corresponding Unicode code points
        let mut chars = self.plain_text.chars();
        let mut segment_index = 0;
        let mut segment_value = 0;
        let mut shift_amount = 0;
        
        while let Some(c) = chars.next() {
            let code_point = c as usize;
            segment_value |= code_point << shift_amount;
            shift_amount += 32;
            if shift_amount == 32 {
                segments[segment_index] = segment_value;
                segment_value = 0;
                shift_amount = 0;
                segment_index += 1;
            }
        }
        
        segments
    }

    fn new(plain_text: String, key: [usize; 4], constants: [usize; 4]) -> LEA {
        let mut lea = LEA {
            plain_text: String::new(),
            key: key,
            key_schedule: [[0; 6];24],
            constants: constants,
            // cipher_text: String::new(),
            // decrypted_text: String::new(),
        };
        
        // Pad the plaintext if its length is less than 128 bits
        let mut padded_plain_text = plain_text;
        while padded_plain_text.len() < 16 {
            padded_plain_text.push('\0'); // Pad with null characters
        }
        
        lea.plain_text = padded_plain_text;
        lea.key_schedule = lea.generate_key_schedule();
        return lea;
    }
}

fn generate_key() -> [usize; 4] {
    let mut rng = OsRng::default(); // Create a new instance of OsRng
    let mut chars: [usize; 4] = [0; 4]; // Initialize the array with zeros

    for i in 0..4 {
        // Generate a random 32-bit value and assign it to the array element
        chars[i] = rng.next_u32() as usize;
    }

    chars
}

fn main() {
    println!("{:?}", generate_key());
    let lea = LEA::new("hello".to_owned(), generate_key(), generate_key());
    println!("{:?}", lea.plain_text);
    println!("{:?}", lea.key);
    println!("{:?}", lea.constants);
    println!("{:?}", lea.generate_key_schedule());
}
