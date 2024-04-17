pub mod prelude {
	pub use crate::lea::cipher::{generic_array::{arr}, BlockDecrypt, BlockEncrypt, NewBlockCipher};
}

mod round_key;

pub use cipher;
use core::mem;
use core::ptr;

use cipher::consts::{U8, U16};
use cipher::generic_array::{typenum::Unsigned, GenericArray};
use cipher::{BlockCipher, BlockDecrypt, BlockEncrypt, NewBlockCipher};

use round_key::{RoundKey, Rk144};

pub type Lea128 = Lea<Rk144>;
pub struct Lea<Rk> where
Rk: RoundKey {
	rk: GenericArray<u32, Rk::RkSize>
}

impl<Rk> BlockCipher for Lea<Rk> where
Rk: RoundKey {
	type BlockSize = U16;
	type ParBlocks = U8;
}

impl<Rk> BlockDecrypt for Lea<Rk> where
Rk: RoundKey {
	fn decrypt_block(&self, block: &mut GenericArray<u8, Self::BlockSize>) {
		decrypt_block::<Rk>(&self.rk, block);
	}
}

impl<Rk> BlockEncrypt for Lea<Rk> where
Rk: RoundKey {
	fn encrypt_block(&self, block: &mut GenericArray<u8, Self::BlockSize>) {
		encrypt_block::<Rk>(&self.rk, block);
	}
}

impl<Rk> NewBlockCipher for Lea<Rk> where
Rk: RoundKey {
	type KeySize = Rk::KeySize;

	fn new(key: &GenericArray<u8, Self::KeySize>) -> Self {
		Self { rk: Rk::generate(key) }
	}
}

fn encrypt_block<Rk>(rk: &GenericArray<u32, Rk::RkSize>, block: &mut GenericArray<u8, <Lea<Rk> as BlockCipher>::BlockSize>) where
Rk: RoundKey {
	let block_orig = block;
	let mut block_copy;

	let block_ptr = block_orig.as_ptr().cast::<[u32; 4]>();
	let block_is_aligned = block_ptr.align_offset(mem::align_of::<[u32; 4]>()) == 0;
	let block = if block_is_aligned {
		unsafe { &mut *block_orig.as_mut_ptr().cast::<[u32; 4]>() }
	} else {
		block_copy = unsafe { block_ptr.read_unaligned() };

		&mut block_copy
	};

	cfg_if::cfg_if! {
		if #[cfg(target_endian = "big")] {
			block[0] = block[0].swap_bytes();
			block[1] = block[1].swap_bytes();
			block[2] = block[2].swap_bytes();
			block[3] = block[3].swap_bytes();
		}
	}

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

	// 28 rounds for 192-bit key
	if <Rk::RkSize as Unsigned>::USIZE / 6 >= 28 {
		block[3] = (block[2] ^ rk[148]).wrapping_add(block[3] ^ rk[149]).rotate_right(3);
		block[2] = (block[1] ^ rk[146]).wrapping_add(block[2] ^ rk[147]).rotate_right(5);
		block[1] = (block[0] ^ rk[144]).wrapping_add(block[1] ^ rk[145]).rotate_left(9);
		block[0] = (block[3] ^ rk[154]).wrapping_add(block[0] ^ rk[155]).rotate_right(3);
		block[3] = (block[2] ^ rk[152]).wrapping_add(block[3] ^ rk[153]).rotate_right(5);
		block[2] = (block[1] ^ rk[150]).wrapping_add(block[2] ^ rk[151]).rotate_left(9);
		block[1] = (block[0] ^ rk[160]).wrapping_add(block[1] ^ rk[161]).rotate_right(3);
		block[0] = (block[3] ^ rk[158]).wrapping_add(block[0] ^ rk[159]).rotate_right(5);
		block[3] = (block[2] ^ rk[156]).wrapping_add(block[3] ^ rk[157]).rotate_left(9);
		block[2] = (block[1] ^ rk[166]).wrapping_add(block[2] ^ rk[167]).rotate_right(3);
		block[1] = (block[0] ^ rk[164]).wrapping_add(block[1] ^ rk[165]).rotate_right(5);
		block[0] = (block[3] ^ rk[162]).wrapping_add(block[0] ^ rk[163]).rotate_left(9);

		// 32 rounds for 256-bit key
		if <Rk::RkSize as Unsigned>::USIZE / 6 >= 32 {
			block[3] = (block[2] ^ rk[172]).wrapping_add(block[3] ^ rk[173]).rotate_right(3);
			block[2] = (block[1] ^ rk[170]).wrapping_add(block[2] ^ rk[171]).rotate_right(5);
			block[1] = (block[0] ^ rk[168]).wrapping_add(block[1] ^ rk[169]).rotate_left(9);
			block[0] = (block[3] ^ rk[178]).wrapping_add(block[0] ^ rk[179]).rotate_right(3);
			block[3] = (block[2] ^ rk[176]).wrapping_add(block[3] ^ rk[177]).rotate_right(5);
			block[2] = (block[1] ^ rk[174]).wrapping_add(block[2] ^ rk[175]).rotate_left(9);
			block[1] = (block[0] ^ rk[184]).wrapping_add(block[1] ^ rk[185]).rotate_right(3);
			block[0] = (block[3] ^ rk[182]).wrapping_add(block[0] ^ rk[183]).rotate_right(5);
			block[3] = (block[2] ^ rk[180]).wrapping_add(block[3] ^ rk[181]).rotate_left(9);
			block[2] = (block[1] ^ rk[190]).wrapping_add(block[2] ^ rk[191]).rotate_right(3);
			block[1] = (block[0] ^ rk[188]).wrapping_add(block[1] ^ rk[189]).rotate_right(5);
			block[0] = (block[3] ^ rk[186]).wrapping_add(block[0] ^ rk[187]).rotate_left(9);
		}
	}

	cfg_if::cfg_if! {
		if #[cfg(target_endian = "big")] {
			block[0] = block[0].swap_bytes();
			block[1] = block[1].swap_bytes();
			block[2] = block[2].swap_bytes();
			block[3] = block[3].swap_bytes();
		}
	}

	if !block_is_aligned {
		unsafe {
			ptr::write_unaligned(block_orig.as_mut_ptr().cast::<[u32; 4]>(), *block);
		}
	}
}

fn decrypt_block<Rk>(rk: &GenericArray<u32, Rk::RkSize>, block: &mut GenericArray<u8, <Lea<Rk> as BlockCipher>::BlockSize>) where
Rk: RoundKey {
	let block_orig = block;
	let mut block_copy;

	let block_ptr = block_orig.as_ptr().cast::<[u32; 4]>();
	let block_is_aligned = block_ptr.align_offset(mem::align_of::<[u32; 4]>()) == 0;
	let block = if block_is_aligned {
		unsafe { &mut *block_orig.as_mut_ptr().cast::<[u32; 4]>() }
	} else {
		block_copy = unsafe { block_ptr.read_unaligned() };

		&mut block_copy
	};

	cfg_if::cfg_if! {
		if #[cfg(target_endian = "big")] {
			block[0] = block[0].swap_bytes();
			block[1] = block[1].swap_bytes();
			block[2] = block[2].swap_bytes();
			block[3] = block[3].swap_bytes();
		}
	}

	// 28 rounds for 192-bit key
	if <Rk::RkSize as Unsigned>::USIZE / 6 >= 28 {
		// 32 rounds for 256-bit key
		if <Rk::RkSize as Unsigned>::USIZE / 6 >= 32 {
			block[0] = block[0].rotate_right(9).wrapping_sub(block[3] ^ rk[186]) ^ rk[187];
			block[1] = block[1].rotate_left(5).wrapping_sub(block[0] ^ rk[188]) ^ rk[189];
			block[2] = block[2].rotate_left(3).wrapping_sub(block[1] ^ rk[190]) ^ rk[191];
			block[3] = block[3].rotate_right(9).wrapping_sub(block[2] ^ rk[180]) ^ rk[181];
			block[0] = block[0].rotate_left(5).wrapping_sub(block[3] ^ rk[182]) ^ rk[183];
			block[1] = block[1].rotate_left(3).wrapping_sub(block[0] ^ rk[184]) ^ rk[185];
			block[2] = block[2].rotate_right(9).wrapping_sub(block[1] ^ rk[174]) ^ rk[175];
			block[3] = block[3].rotate_left(5).wrapping_sub(block[2] ^ rk[176]) ^ rk[177];
			block[0] = block[0].rotate_left(3).wrapping_sub(block[3] ^ rk[178]) ^ rk[179];
			block[1] = block[1].rotate_right(9).wrapping_sub(block[0] ^ rk[168]) ^ rk[169];
			block[2] = block[2].rotate_left(5).wrapping_sub(block[1] ^ rk[170]) ^ rk[171];
			block[3] = block[3].rotate_left(3).wrapping_sub(block[2] ^ rk[172]) ^ rk[173];
		}

		block[0] = block[0].rotate_right(9).wrapping_sub(block[3] ^ rk[162]) ^ rk[163];
		block[1] = block[1].rotate_left(5).wrapping_sub(block[0] ^ rk[164]) ^ rk[165];
		block[2] = block[2].rotate_left(3).wrapping_sub(block[1] ^ rk[166]) ^ rk[167];
		block[3] = block[3].rotate_right(9).wrapping_sub(block[2] ^ rk[156]) ^ rk[157];
		block[0] = block[0].rotate_left(5).wrapping_sub(block[3] ^ rk[158]) ^ rk[159];
		block[1] = block[1].rotate_left(3).wrapping_sub(block[0] ^ rk[160]) ^ rk[161];
		block[2] = block[2].rotate_right(9).wrapping_sub(block[1] ^ rk[150]) ^ rk[151];
		block[3] = block[3].rotate_left(5).wrapping_sub(block[2] ^ rk[152]) ^ rk[153];
		block[0] = block[0].rotate_left(3).wrapping_sub(block[3] ^ rk[154]) ^ rk[155];
		block[1] = block[1].rotate_right(9).wrapping_sub(block[0] ^ rk[144]) ^ rk[145];
		block[2] = block[2].rotate_left(5).wrapping_sub(block[1] ^ rk[146]) ^ rk[147];
		block[3] = block[3].rotate_left(3).wrapping_sub(block[2] ^ rk[148]) ^ rk[149];
	}

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

	cfg_if::cfg_if! {
		if #[cfg(target_endian = "big")] {
			block[0] = block[0].swap_bytes();
			block[1] = block[1].swap_bytes();
			block[2] = block[2].swap_bytes();
			block[3] = block[3].swap_bytes();
		}
	}

	if !block_is_aligned {
		unsafe {
			ptr::write_unaligned(block_orig.as_mut_ptr().cast::<[u32; 4]>(), *block);
		}
	}
}
