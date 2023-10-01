# sutils
Personal socket tools for go

# This module works for a certain socket body structure:
- section
```
section_content_length (8 bytes, base-256) + section_content
```
- body:
```
section + section + ...
```

# Note:
customize base-256 with 8 bytes, supports content length lower than   
***15.506782791695285 PB***  

feel free to do anything
