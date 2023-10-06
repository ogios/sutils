# sutils
Personal socket tools for go

# This module works for a certain socket body structure:
- section
```
section_content_length (base-255) + section_content
```
- body:
```
section + section + ...
```

# Note:
~~customize base-256 with 8 bytes, supports content length lower than~~   
~~***15.506782791695285 PB***~~  

Changed to base*255 with the last byte of `255` ends the length block  
which means less memory usage and smaller request body  
and no length limitation!

feel free to do anything
