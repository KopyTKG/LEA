use core::mem;

// Constants for key expansion
const CONSTVAL: [u32; 8] = [0xC3EFE9DB, 0x44626B02, 0x79E27C8A, 0x78DF30EC, 0x715EA49E, 0xC785DA0A, 0xE04EF22A, 0xE5C40957];

pub	fn generate(key: [u32; 16] ) -> [u32; 144] {
    // Convert key to u32 array
    let key_ptr = key.as_ptr().cast::<[u32; 4]>();

    // Check if key is aligned
    let key_is_aligned = key_ptr.align_offset(mem::align_of::<[u32; 4]>()) == 0;

    // Read key
    let mut rk_t = if key_is_aligned {
        unsafe { key_ptr.read() }
    } else {
        unsafe { key_ptr.read_unaligned() }
    };

    // Initialize round key
    let mut rk : [u32; 144] = [0; 144];

    // Copy key to round key
    for i in 0..4 {
        rk[i] = rk_t[i];
    }

    // Generate round key
    for i in 0..24 {
        let t0 = CONSTVAL[i % 4].rotate_left(i as u32);
        let t1 = t0.rotate_left(1);
        let t2 = t1.rotate_left(1);
        let t3 = t2.rotate_left(1);
        rk_t[0] = rk_t[0].wrapping_add(t0).rotate_left(1);
        rk_t[1] = rk_t[1].wrapping_add(t1).rotate_left(3);
        rk_t[2] = rk_t[2].wrapping_add(t2).rotate_left(6);
        rk_t[3] = rk_t[3].wrapping_add(t3).rotate_left(11);

        rk[6 * i + 0] = rk_t[0];
        rk[6 * i + 1] = rk_t[1];
        rk[6 * i + 2] = rk_t[2];
        rk[6 * i + 3] = rk_t[1];
        rk[6 * i + 4] = rk_t[3];
        rk[6 * i + 5] = rk_t[1];
    }

    rk
}

