use std::io::{Read, Write};
use std::fs::File;
use std::path::Path;

mod lea;
use lea::{Lea,Lea128};
mod round_key;


fn create_blocks(buffer: Vec<u8>) -> Vec<[u8; 16]> {
    let mut blocks: Vec<[u8;16]> = Vec::new();
    let full_blocks = buffer.len() / 16;
    let remaining_bytes = buffer.len() % 16;

    let total_blocks = if remaining_bytes == 0 {
        full_blocks
    } else {
        full_blocks + 1
    };

    for i in 0..total_blocks {
        let mut block: [u8;16] = [0x20; 16];
        for j in 0..16 {
            let index = i * 16 + j;
            if index < buffer.len() {
                block[j] = buffer[index];
            }
        }
        blocks.push(block);
    }
    blocks
}


fn main() -> std::io::Result<()> {
    // pre-defined key
    let key: [u32; 16] = [0x0F, 0x1E, 0x2D, 0x3C, 0x4B, 0x5A, 0x69, 0x78, 0x87, 0x96, 0xA5, 0xB4, 0xC3, 0xD2, 0xE1, 0xF0];
    // create a new instance of the LEA block cipher
    let lea: Lea = Lea::new(&key);

    let path = Path::new("source.txt");
    // Try to open the file or create it if it does not exist
    let mut file = if path.exists() {
        File::open(path)?
    } else {
        File::create(path)?
    };
    // Read the file into a buffer
    let mut buffer = Vec::new();
    file.read_to_end(&mut buffer)?;
    let mut blocks = create_blocks(buffer);

    // Encrypt the blocks
    let mut save = File::create("encrypted.txt")?;
    for bl in blocks.iter_mut() {
        let mut block = bl.clone();
        lea.encrypt_block(&mut block);
        save.write_all(&block)?;
    }
    
    // Read the encrypted file
    let mut enc = File::open("encrypted.txt")?;
    let mut buffer = Vec::new();
    enc.read_to_end(&mut buffer)?;
    let mut blocks = create_blocks(buffer); 
    
    // Decrypt the blocks
    let mut save = File::create("decrypted.txt")?;
    for bl in blocks.iter_mut() {
        let mut block = bl.clone();
        lea.decrypt_block(&mut block);
        save.write_all(&block)?;
    }
    
    Ok(())
}
