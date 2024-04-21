use core::mem;
use crate::round_key;


pub trait Lea128 {
    fn encrypt_block(&self, block: &mut [u8; 16]);
    fn decrypt_block(&self, block: &mut [u8; 16]);
    fn new(key: &[u32; 16]) -> Self where Self: Sized;
}

pub struct Lea {
    rk: [u32; 144],
    key: [u32; 16],
}

impl Lea128 for Lea {
    fn new(key: &[u32; 16]) -> Self {
        let mut lea = Self {
            rk: [0; 144],
            key: *key,
        };
        lea.rk = round_key::generate(lea.key);
        lea
    }


    fn encrypt_block(&self, block: &mut [u8; 16]) {
        let block_orig = block;
        let rk = self.rk;
        let mut block_copy;

        let block_ptr = block_orig.as_ptr().cast::<[u32; 4]>();
        let block_is_aligned = block_ptr.align_offset(mem::align_of::<[u32; 4]>()) == 0;
        let block = if block_is_aligned {
            unsafe { &mut *block_orig.as_mut_ptr().cast::<[u32; 4]>() }
        } else {
            block_copy = unsafe { block_ptr.read_unaligned() };

            &mut block_copy
        };

        // 24 rounds for 128-bit key
        block[3] = (block[2] ^ rk[4]).wrapping_add(block[3] ^ rk[5]).rotate_right(3);
        block[2] = (block[1] ^ rk[2]).wrapping_add(block[2] ^ rk[3]).rotate_right(5);
        block[1] = (block[0] ^ rk[0]).wrapping_add(block[1] ^ rk[1]).rotate_left(9);
        block[0] = (block[3] ^ rk[10]).wrapping_add(block[0] ^ rk[11]).rotate_right(3);
        block[3] = (block[2] ^ rk[8]).wrapping_add(block[3] ^ rk[9]).rotate_right(5);
        block[2] = (block[1] ^ rk[6]).wrapping_add(block[2] ^ rk[7]).rotate_left(9);
        block[1] = (block[0] ^ rk[16]).wrapping_add(block[1] ^ rk[17]).rotate_right(3);
        block[0] = (block[3] ^ rk[14]).wrapping_add(block[0] ^ rk[15]).rotate_right(5);
        block[3] = (block[2] ^ rk[12]).wrapping_add(block[3] ^ rk[13]).rotate_left(9);
        block[2] = (block[1] ^ rk[22]).wrapping_add(block[2] ^ rk[23]).rotate_right(3);
        block[1] = (block[0] ^ rk[20]).wrapping_add(block[1] ^ rk[21]).rotate_right(5);
        block[0] = (block[3] ^ rk[18]).wrapping_add(block[0] ^ rk[19]).rotate_left(9);

        block[3] = (block[2] ^ rk[28]).wrapping_add(block[3] ^ rk[29]).rotate_right(3);
        block[2] = (block[1] ^ rk[26]).wrapping_add(block[2] ^ rk[27]).rotate_right(5);
        block[1] = (block[0] ^ rk[24]).wrapping_add(block[1] ^ rk[25]).rotate_left(9);
        block[0] = (block[3] ^ rk[34]).wrapping_add(block[0] ^ rk[35]).rotate_right(3);
        block[3] = (block[2] ^ rk[32]).wrapping_add(block[3] ^ rk[33]).rotate_right(5);
        block[2] = (block[1] ^ rk[30]).wrapping_add(block[2] ^ rk[31]).rotate_left(9);
        block[1] = (block[0] ^ rk[40]).wrapping_add(block[1] ^ rk[41]).rotate_right(3);
        block[0] = (block[3] ^ rk[38]).wrapping_add(block[0] ^ rk[39]).rotate_right(5);
        block[3] = (block[2] ^ rk[36]).wrapping_add(block[3] ^ rk[37]).rotate_left(9);
        block[2] = (block[1] ^ rk[46]).wrapping_add(block[2] ^ rk[47]).rotate_right(3);
        block[1] = (block[0] ^ rk[44]).wrapping_add(block[1] ^ rk[45]).rotate_right(5);
        block[0] = (block[3] ^ rk[42]).wrapping_add(block[0] ^ rk[43]).rotate_left(9);

        block[3] = (block[2] ^ rk[52]).wrapping_add(block[3] ^ rk[53]).rotate_right(3);
        block[2] = (block[1] ^ rk[50]).wrapping_add(block[2] ^ rk[51]).rotate_right(5);
        block[1] = (block[0] ^ rk[48]).wrapping_add(block[1] ^ rk[49]).rotate_left(9);
        block[0] = (block[3] ^ rk[58]).wrapping_add(block[0] ^ rk[59]).rotate_right(3);
        block[3] = (block[2] ^ rk[56]).wrapping_add(block[3] ^ rk[57]).rotate_right(5);
        block[2] = (block[1] ^ rk[54]).wrapping_add(block[2] ^ rk[55]).rotate_left(9);
        block[1] = (block[0] ^ rk[64]).wrapping_add(block[1] ^ rk[65]).rotate_right(3);
        block[0] = (block[3] ^ rk[62]).wrapping_add(block[0] ^ rk[63]).rotate_right(5);
        block[3] = (block[2] ^ rk[60]).wrapping_add(block[3] ^ rk[61]).rotate_left(9);
        block[2] = (block[1] ^ rk[70]).wrapping_add(block[2] ^ rk[71]).rotate_right(3);
        block[1] = (block[0] ^ rk[68]).wrapping_add(block[1] ^ rk[69]).rotate_right(5);
        block[0] = (block[3] ^ rk[66]).wrapping_add(block[0] ^ rk[67]).rotate_left(9);

        block[3] = (block[2] ^ rk[76]).wrapping_add(block[3] ^ rk[77]).rotate_right(3);
        block[2] = (block[1] ^ rk[74]).wrapping_add(block[2] ^ rk[75]).rotate_right(5);
        block[1] = (block[0] ^ rk[72]).wrapping_add(block[1] ^ rk[73]).rotate_left(9);
        block[0] = (block[3] ^ rk[82]).wrapping_add(block[0] ^ rk[83]).rotate_right(3);
        block[3] = (block[2] ^ rk[80]).wrapping_add(block[3] ^ rk[81]).rotate_right(5);
        block[2] = (block[1] ^ rk[78]).wrapping_add(block[2] ^ rk[79]).rotate_left(9);
        block[1] = (block[0] ^ rk[88]).wrapping_add(block[1] ^ rk[89]).rotate_right(3);
        block[0] = (block[3] ^ rk[86]).wrapping_add(block[0] ^ rk[87]).rotate_right(5);
        block[3] = (block[2] ^ rk[84]).wrapping_add(block[3] ^ rk[85]).rotate_left(9);
        block[2] = (block[1] ^ rk[94]).wrapping_add(block[2] ^ rk[95]).rotate_right(3);
        block[1] = (block[0] ^ rk[92]).wrapping_add(block[1] ^ rk[93]).rotate_right(5);
        block[0] = (block[3] ^ rk[90]).wrapping_add(block[0] ^ rk[91]).rotate_left(9);

        block[3] = (block[2] ^ rk[100]).wrapping_add(block[3] ^ rk[101]).rotate_right(3);
        block[2] = (block[1] ^ rk[98]).wrapping_add(block[2] ^ rk[99]).rotate_right(5);
        block[1] = (block[0] ^ rk[96]).wrapping_add(block[1] ^ rk[97]).rotate_left(9);
        block[0] = (block[3] ^ rk[106]).wrapping_add(block[0] ^ rk[107]).rotate_right(3);
        block[3] = (block[2] ^ rk[104]).wrapping_add(block[3] ^ rk[105]).rotate_right(5);
        block[2] = (block[1] ^ rk[102]).wrapping_add(block[2] ^ rk[103]).rotate_left(9);
        block[1] = (block[0] ^ rk[112]).wrapping_add(block[1] ^ rk[113]).rotate_right(3);
        block[0] = (block[3] ^ rk[110]).wrapping_add(block[0] ^ rk[111]).rotate_right(5);
        block[3] = (block[2] ^ rk[108]).wrapping_add(block[3] ^ rk[109]).rotate_left(9);
        block[2] = (block[1] ^ rk[118]).wrapping_add(block[2] ^ rk[119]).rotate_right(3);
        block[1] = (block[0] ^ rk[116]).wrapping_add(block[1] ^ rk[117]).rotate_right(5);
        block[0] = (block[3] ^ rk[114]).wrapping_add(block[0] ^ rk[115]).rotate_left(9);

        block[3] = (block[2] ^ rk[124]).wrapping_add(block[3] ^ rk[125]).rotate_right(3);
        block[2] = (block[1] ^ rk[122]).wrapping_add(block[2] ^ rk[123]).rotate_right(5);
        block[1] = (block[0] ^ rk[120]).wrapping_add(block[1] ^ rk[121]).rotate_left(9);
        block[0] = (block[3] ^ rk[130]).wrapping_add(block[0] ^ rk[131]).rotate_right(3);
        block[3] = (block[2] ^ rk[128]).wrapping_add(block[3] ^ rk[129]).rotate_right(5);
        block[2] = (block[1] ^ rk[126]).wrapping_add(block[2] ^ rk[127]).rotate_left(9);
        block[1] = (block[0] ^ rk[136]).wrapping_add(block[1] ^ rk[137]).rotate_right(3);
        block[0] = (block[3] ^ rk[134]).wrapping_add(block[0] ^ rk[135]).rotate_right(5);
        block[3] = (block[2] ^ rk[132]).wrapping_add(block[3] ^ rk[133]).rotate_left(9);
        block[2] = (block[1] ^ rk[142]).wrapping_add(block[2] ^ rk[143]).rotate_right(3);
        block[1] = (block[0] ^ rk[140]).wrapping_add(block[1] ^ rk[141]).rotate_right(5);
        block[0] = (block[3] ^ rk[138]).wrapping_add(block[0] ^ rk[139]).rotate_left(9);
    }

    fn decrypt_block(&self, block: &mut [u8; 16]) {
        let block_orig = block;
        let rk = self.rk;
        let mut block_copy;

        let block_ptr = block_orig.as_ptr().cast::<[u32; 4]>();
        let block_is_aligned = block_ptr.align_offset(mem::align_of::<[u32; 4]>()) == 0;
        let block = if block_is_aligned {
            unsafe { &mut *block_orig.as_mut_ptr().cast::<[u32; 4]>() }
        } else {
            block_copy = unsafe { block_ptr.read_unaligned() };

            &mut block_copy
        };

        // 24 rounds for 128-bit key
        block[0] = block[0].rotate_right(9).wrapping_sub(block[3] ^ rk[138]) ^ rk[139];
        block[1] = block[1].rotate_left(5).wrapping_sub(block[0] ^ rk[140]) ^ rk[141];
        block[2] = block[2].rotate_left(3).wrapping_sub(block[1] ^ rk[142]) ^ rk[143];
        block[3] = block[3].rotate_right(9).wrapping_sub(block[2] ^ rk[132]) ^ rk[133];
        block[0] = block[0].rotate_left(5).wrapping_sub(block[3] ^ rk[134]) ^ rk[135];
        block[1] = block[1].rotate_left(3).wrapping_sub(block[0] ^ rk[136]) ^ rk[137];
        block[2] = block[2].rotate_right(9).wrapping_sub(block[1] ^ rk[126]) ^ rk[127];
        block[3] = block[3].rotate_left(5).wrapping_sub(block[2] ^ rk[128]) ^ rk[129];
        block[0] = block[0].rotate_left(3).wrapping_sub(block[3] ^ rk[130]) ^ rk[131];
        block[1] = block[1].rotate_right(9).wrapping_sub(block[0] ^ rk[120]) ^ rk[121];
        block[2] = block[2].rotate_left(5).wrapping_sub(block[1] ^ rk[122]) ^ rk[123];
        block[3] = block[3].rotate_left(3).wrapping_sub(block[2] ^ rk[124]) ^ rk[125];

        block[0] = block[0].rotate_right(9).wrapping_sub(block[3] ^ rk[114]) ^ rk[115];
        block[1] = block[1].rotate_left(5).wrapping_sub(block[0] ^ rk[116]) ^ rk[117];
        block[2] = block[2].rotate_left(3).wrapping_sub(block[1] ^ rk[118]) ^ rk[119];
        block[3] = block[3].rotate_right(9).wrapping_sub(block[2] ^ rk[108]) ^ rk[109];
        block[0] = block[0].rotate_left(5).wrapping_sub(block[3] ^ rk[110]) ^ rk[111];
        block[1] = block[1].rotate_left(3).wrapping_sub(block[0] ^ rk[112]) ^ rk[113];
        block[2] = block[2].rotate_right(9).wrapping_sub(block[1] ^ rk[102]) ^ rk[103];
        block[3] = block[3].rotate_left(5).wrapping_sub(block[2] ^ rk[104]) ^ rk[105];
        block[0] = block[0].rotate_left(3).wrapping_sub(block[3] ^ rk[106]) ^ rk[107];
        block[1] = block[1].rotate_right(9).wrapping_sub(block[0] ^ rk[96]) ^ rk[97];
        block[2] = block[2].rotate_left(5).wrapping_sub(block[1] ^ rk[98]) ^ rk[99];
        block[3] = block[3].rotate_left(3).wrapping_sub(block[2] ^ rk[100]) ^ rk[101];

        block[0] = block[0].rotate_right(9).wrapping_sub(block[3] ^ rk[90]) ^ rk[91];
        block[1] = block[1].rotate_left(5).wrapping_sub(block[0] ^ rk[92]) ^ rk[93];
        block[2] = block[2].rotate_left(3).wrapping_sub(block[1] ^ rk[94]) ^ rk[95];
        block[3] = block[3].rotate_right(9).wrapping_sub(block[2] ^ rk[84]) ^ rk[85];
        block[0] = block[0].rotate_left(5).wrapping_sub(block[3] ^ rk[86]) ^ rk[87];
        block[1] = block[1].rotate_left(3).wrapping_sub(block[0] ^ rk[88]) ^ rk[89];
        block[2] = block[2].rotate_right(9).wrapping_sub(block[1] ^ rk[78]) ^ rk[79];
        block[3] = block[3].rotate_left(5).wrapping_sub(block[2] ^ rk[80]) ^ rk[81];
        block[0] = block[0].rotate_left(3).wrapping_sub(block[3] ^ rk[82]) ^ rk[83];
        block[1] = block[1].rotate_right(9).wrapping_sub(block[0] ^ rk[72]) ^ rk[73];
        block[2] = block[2].rotate_left(5).wrapping_sub(block[1] ^ rk[74]) ^ rk[75];
        block[3] = block[3].rotate_left(3).wrapping_sub(block[2] ^ rk[76]) ^ rk[77];

        block[0] = block[0].rotate_right(9).wrapping_sub(block[3] ^ rk[66]) ^ rk[67];
        block[1] = block[1].rotate_left(5).wrapping_sub(block[0] ^ rk[68]) ^ rk[69];
        block[2] = block[2].rotate_left(3).wrapping_sub(block[1] ^ rk[70]) ^ rk[71];
        block[3] = block[3].rotate_right(9).wrapping_sub(block[2] ^ rk[60]) ^ rk[61];
        block[0] = block[0].rotate_left(5).wrapping_sub(block[3] ^ rk[62]) ^ rk[63];
        block[1] = block[1].rotate_left(3).wrapping_sub(block[0] ^ rk[64]) ^ rk[65];
        block[2] = block[2].rotate_right(9).wrapping_sub(block[1] ^ rk[54]) ^ rk[55];
        block[3] = block[3].rotate_left(5).wrapping_sub(block[2] ^ rk[56]) ^ rk[57];
        block[0] = block[0].rotate_left(3).wrapping_sub(block[3] ^ rk[58]) ^ rk[59];
        block[1] = block[1].rotate_right(9).wrapping_sub(block[0] ^ rk[48]) ^ rk[49];
        block[2] = block[2].rotate_left(5).wrapping_sub(block[1] ^ rk[50]) ^ rk[51];
        block[3] = block[3].rotate_left(3).wrapping_sub(block[2] ^ rk[52]) ^ rk[53];

        block[0] = block[0].rotate_right(9).wrapping_sub(block[3] ^ rk[42]) ^ rk[43];
        block[1] = block[1].rotate_left(5).wrapping_sub(block[0] ^ rk[44]) ^ rk[45];
        block[2] = block[2].rotate_left(3).wrapping_sub(block[1] ^ rk[46]) ^ rk[47];
        block[3] = block[3].rotate_right(9).wrapping_sub(block[2] ^ rk[36]) ^ rk[37];
        block[0] = block[0].rotate_left(5).wrapping_sub(block[3] ^ rk[38]) ^ rk[39];
        block[1] = block[1].rotate_left(3).wrapping_sub(block[0] ^ rk[40]) ^ rk[41];
        block[2] = block[2].rotate_right(9).wrapping_sub(block[1] ^ rk[30]) ^ rk[31];
        block[3] = block[3].rotate_left(5).wrapping_sub(block[2] ^ rk[32]) ^ rk[33];
        block[0] = block[0].rotate_left(3).wrapping_sub(block[3] ^ rk[34]) ^ rk[35];
        block[1] = block[1].rotate_right(9).wrapping_sub(block[0] ^ rk[24]) ^ rk[25];
        block[2] = block[2].rotate_left(5).wrapping_sub(block[1] ^ rk[26]) ^ rk[27];
        block[3] = block[3].rotate_left(3).wrapping_sub(block[2] ^ rk[28]) ^ rk[29];

        block[0] = block[0].rotate_right(9).wrapping_sub(block[3] ^ rk[18]) ^ rk[19];
        block[1] = block[1].rotate_left(5).wrapping_sub(block[0] ^ rk[20]) ^ rk[21];
        block[2] = block[2].rotate_left(3).wrapping_sub(block[1] ^ rk[22]) ^ rk[23];
        block[3] = block[3].rotate_right(9).wrapping_sub(block[2] ^ rk[12]) ^ rk[13];
        block[0] = block[0].rotate_left(5).wrapping_sub(block[3] ^ rk[14]) ^ rk[15];
        block[1] = block[1].rotate_left(3).wrapping_sub(block[0] ^ rk[16]) ^ rk[17];
        block[2] = block[2].rotate_right(9).wrapping_sub(block[1] ^ rk[6]) ^ rk[7];
        block[3] = block[3].rotate_left(5).wrapping_sub(block[2] ^ rk[8]) ^ rk[9];
        block[0] = block[0].rotate_left(3).wrapping_sub(block[3] ^ rk[10]) ^ rk[11];
        block[1] = block[1].rotate_right(9).wrapping_sub(block[0] ^ rk[0]) ^ rk[1];
        block[2] = block[2].rotate_left(5).wrapping_sub(block[1] ^ rk[2]) ^ rk[3];
        block[3] = block[3].rotate_left(3).wrapping_sub(block[2] ^ rk[4]) ^ rk[5];
    }
}
