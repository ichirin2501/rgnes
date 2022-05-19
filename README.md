# rgnes

[![Build Status](https://github.com/ichirin2501/rgnes/workflows/Test/badge.svg?branch=master)](https://github.com/ichirin2501/rgnes/actions)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

This NES emulator is implemented for my study

Most of the code in rgnes was based on the following.
It's great because it's written so simple.

- https://github.com/fogleman/nes

## TODOs

- decay ppu openbus
- apu emulation
  - (I'm not interested in sound so I probably won't implement it...)

.. and tons more


## Tests

| Test | SingleRom | Result |
| - | - | - |
| blargg_ppu_tests_2005.09.15b | palette_ram.nes | OK |
| blargg_ppu_tests_2005.09.15b | sprite_ram.nes  | OK |
| blargg_ppu_tests_2005.09.15b | vbl_clear_time.nes | OK |
| blargg_ppu_tests_2005.09.15b | vram_access.nes  | OK |
| branch_timing_tests | 1.Branch_Basics.nes | OK |
| branch_timing_tests | 2.Backward_Branch.nes | OK |
| branch_timing_tests | 3.Forward_Branch.nes  | OK |
| cpu_dummy_reads | cpu_dummy_reads.nes | OK |
| cpu_dummy_writes | cpu_dummy_writes_oam.nes | OK |
| cpu_dummy_writes | cpu_dummy_writes_ppumem.nes | OK |
| cpu_exec_space | test_cpu_exec_space_apu.nes | Failed |
| cpu_exec_space | test_cpu_exec_space_ppuio.nes | OK |
| cpu_timing_test6 | cpu_timing_test.nes | OK |
| instr_misc | 01-abs_x_wrap.nes | OK |
| instr_misc | 02-branch_wrap.nes | OK |
| instr_misc | 03-dummy_reads.nes | OK |
| instr_misc | 04-dummy_reads_apu.nes | Failed |
| instr_test-v5 | 01-basics.nes | OK |
| instr_test-v5 | 02-implied.nes | OK |
| instr_test-v5 | 03-immediate.nes | OK |
| instr_test-v5 | 04-zero_page.nes | OK |
| instr_test-v5 | 05-zp_xy.nes | OK |
| instr_test-v5 | 06-absolute.nes | OK |
| instr_test-v5 | 07-abs_xy.nes | OK |
| instr_test-v5 | 08-ind_x.nes | OK |
| instr_test-v5 | 09-ind_y.nes | OK |
| instr_test-v5 | 10-branches.nes | OK |
| instr_test-v5 | 11-stack.nes | OK |
| instr_test-v5 | 12-jmp_jsr.nes | OK |
| instr_test-v5 | 13-rts.nes | OK |
| instr_test-v5 | 14-rti.nes | OK |
| instr_test-v5 | 15-brk.nes | OK |
| instr_test-v5 | 16-special.nes | OK |
| nestest | nestest.nes | OK |
| oam_read | oam_read.nes | OK |
| oam_stress | oam_stress.nes | OK |
| ppu_open_bus | ppu_open_bus.nes | Failed |
| ppu_read_buffer | test_ppu_read_buffer.nes | OK |
| ppu_vbl_nmi | 01-vbl_basics.nes | OK |
| ppu_vbl_nmi | 02-vbl_set_time.nes | OK |
| ppu_vbl_nmi | 03-vbl_clear_time.nes | OK |
| ppu_vbl_nmi | 04-nmi_control.nes  | OK |
| ppu_vbl_nmi | 05-nmi_timing.nes  | OK |
| ppu_vbl_nmi | 06-suppression.nes | OK |
| ppu_vbl_nmi | 07-nmi_on_timing.nes | OK |
| ppu_vbl_nmi | 08-nmi_off_timing.nes | OK |
| ppu_vbl_nmi | 09-even_odd_frames.nes | OK |
| ppu_vbl_nmi | 10-even_odd_timing.nes | Failed |
| sprite_hit_tests_2005.10.05 | 01.basics.nes | OK |
| sprite_hit_tests_2005.10.05 | 02.alignment.nes | OK |
| sprite_hit_tests_2005.10.05 | 03.corners.nes | OK |
| sprite_hit_tests_2005.10.05 | 04.flip.nes | OK |
| sprite_hit_tests_2005.10.05 | 05.left_clip.nes | OK |
| sprite_hit_tests_2005.10.05 | 06.right_edge.nes | OK |
| sprite_hit_tests_2005.10.05 | 07.screen_bottom.nes | OK |
| sprite_hit_tests_2005.10.05 | 08.double_height.nes | OK |
| sprite_hit_tests_2005.10.05 | 09.timing_basics.nes | OK |
| sprite_hit_tests_2005.10.05 | 10.timing_order.nes | OK |
| sprite_hit_tests_2005.10.05 | 11.edge_timing.nes | OK |
| sprite_overflow_tests | 1.Basics.nes | OK |
| sprite_overflow_tests | 2.Details.nes | OK |
| sprite_overflow_tests | 3.Timing.nes | Failed |
| sprite_overflow_tests | 4.Obscure.nes | Failed |
| sprite_overflow_tests | 5.Emulator.nes | OK |
| vbl_nmi_timing | 1.frame_basics.nes | OK |
| vbl_nmi_timing | 2.vbl_timing.nes | OK |
| vbl_nmi_timing | 3.even_odd_frames.nes | OK |
| vbl_nmi_timing | 4.vbl_clear_timing.nes | OK |
| vbl_nmi_timing | 5.nmi_suppression.nes | OK |
| vbl_nmi_timing | 6.nmi_disable.nes | OK |
| vbl_nmi_timing | 7.nmi_timing.nes | OK |

https://github.com/christopherpow/nes-test-roms/blob/master/other/RasterDemo.NES  
![RasterDemo.nes](/images/RasterDemo.gif)  

https://github.com/christopherpow/nes-test-roms/blob/master/other/RasterTest1.NES  
![RasterTest1.nes](/images/RasterTest1.gif)  

https://www.nesdev.org/wiki/Emulator_tests#PPU_Tests  
scanline.nes (This test ROM was created by Quietust)  
![scanline.nes](/images/scanline.gif)  

## References

- https://www.nesdev.org/wiki/NES_reference_guide
- https://github.com/fogleman/nes
- [ファミコンエミュレータの創り方　- Hello, World!編 -](https://qiita.com/bokuweb/items/1575337bef44ae82f4d3)
- [Writing NES Emulator in Rust](https://bugzmanov.github.io/nes_ebook/chapter_1.html)
- [ｷﾞｺ猫でもわかるファミコンプログラミング](http://gikofami.fc2web.com/index.html)

## License
rgnes is licensed under the MIT license
