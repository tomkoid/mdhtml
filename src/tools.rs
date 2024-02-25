use std::path::Path;

use crate::args::Convert;

// pub fn get_filename(args: &Convert) -> String {
//     let mut true_output_file = String::new();
//
//     if args.output == "[INPUT].html" {
//         // I HAVE 0 RUST EXPERIENCE LOL, please forgive me
//
//         // get rid of the whole file extension and add .html
//         let output_file = Arc::new(Mutex::new(args.input.clone()));
//         let extension_output_file = Arc::clone(&output_file).lock().unwrap().clone();
//
//         let extension = Path::new(&extension_output_file).extension();
//         if extension.is_some() {
//             true_output_file = output_file.lock().unwrap().replace(
//                 &format!(".{}", extension.unwrap().to_str().unwrap()),
//                 ".html",
//             );
//         } else {
//             true_output_file = format!("{}.html", output_file.lock().unwrap().clone());
//         }
//     }
//
//     true_output_file
// }

pub fn get_filename(args: &Convert) -> String {
    let true_output_file: String;

    if args.output == "[INPUT].html" {
        let input_file = Path::new(&args.input);
        let stem = input_file.file_stem().unwrap().to_str().unwrap();
        true_output_file = format!("{}.html", stem);
    } else {
        true_output_file = args.output.replace("[INPUT]", &args.input);
    }

    true_output_file
}
