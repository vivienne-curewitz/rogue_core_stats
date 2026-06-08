// uesave_bridge.h
#include <stddef.h>
#include <stdint.h>

void free_rust_string(char* ptr);
char* convert_to_json(const uint8_t* file_bytes, size_t len);
