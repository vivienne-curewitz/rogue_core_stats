use std::ffi::{CStr, CString};
use std::os::raw::c_char;

// Safety wrapper to free strings allocated by Rust from the Go side
#[no_mangle]
pub unsafe extern "C" fn free_rust_string(ptr: *mut c_char) {
    if !ptr.is_null() {
        let _ = CString::from_raw(ptr);
    }
}

#[no_mangle]
pub unsafe extern "C" fn convert_to_json(file_bytes: *const u8, len: usize) -> *mut c_char {
    if file_bytes.is_null() {
        return std::ptr::null_mut();
    }

    // Convert the raw pointer into a Rust slice
    let bytes = std::slice::from_raw_parts(file_bytes, len);

    // Create a cursor and read using uesave
    let mut cursor = std::io::Cursor::new(bytes);

    // Perform uesave read operation
    match uesave::Save::read(&mut cursor) {
        Ok(save) => {
            // Convert to a JSON string (uesave supports serde serialization)
            if let Ok(json_str) = serde_json::to_string(&save) {
                if let Ok(c_str) = CString::new(json_str) {
                    return c_str.into_raw(); // Hand over memory ownership to Go/C
                }
            }
            std::ptr::null_mut()
        }
        Err(_) => std::ptr::null_mut(),
    }
}
