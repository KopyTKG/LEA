//std libs
use std::io::prelude::*;
//use std::io::{Read, Write};
use std::fs::{File};
use std::path::Path;

mod lea;
use lea::{prelude::*, Lea128};

fn string_to_u8_array(s: &str) -> [u8; 16] {
    let mut arr = [0u8; 16];
    for (i, c) in s.bytes().enumerate() {
        if i >= 16 {
            break;
        }
        arr[i] = c;
    }
    arr

}


fn u8_array_to_string(arr: &[u8; 16]) -> String {
    let mut s = String::new();
    for c in arr.iter() {
        s.push(*c as char);
    }
    s
}

fn main() -> std::io::Result<()> {
    let path = Path::new("source.txt");

    if !path.exists() {
        File::create("source.txt");
        return Err(std::io::Error::new(std::io::ErrorKind::NotFound, "source.txt not found"));
    }

    let mut file = File::open("source.txt")?;
    let mut buffer = Vec::new();
    file.read_to_end(&mut buffer);

    // split the buffer into u8 array with 16 bytes
    let mut blocks: arr! = arr![u8;16];
    for i in 0..buffer.len() / 16 {
        let mut block: [u8:16] = [0; 16];
        for j in 0..16 {
            block[j] = buffer[i * 16 + j];
        }
        blocks.push(block);
    }
    // pregenerated key for encryption
    let key = arr![u8; 0x0F, 0x1E, 0x2D, 0x3C, 0x4B, 0x5A, 0x69, 0x78, 0x87, 0x96, 0xA5, 0xB4, 0xC3, 0xD2, 0xE1, 0xF0];
    
    let mut file = File::create("encrypted.txt");

    for block in blocks.iter_mut() {
        let mut block = block.clone();
        let lea128 = Lea128::new(&key);
        lea128.encrypt_block(block);
        file.write_all(&block);
    }

    //let mut file = File::create("encrypted.txt");
    //for block in blocks.iter() {
    //    file.write_all(block);
    //}

    Ok(())
    // decryption process
    // lea128.decrypt_block((&mut block).into());
    // let text = u8_array_to_string(&block);
    // println!("{:?}", text);

}
