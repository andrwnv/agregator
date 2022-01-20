import { UnsupportedMediaTypeException } from '@nestjs/common';
import { extname } from 'path';

export const imageFilter = (_, file, callback) => {
    if (!file.originalname.match(/\.(jpg|jpeg|png|gif)$/))
        return callback(new UnsupportedMediaTypeException({},
            "Incorrect file type. Supported JPG, JPEG, PNG, GIF"), false);

    callback(null, true);
}

export const editFileName = (req, file, callback) => {
    const extension = extname(file.originalname).split('.')[1];
    callback(null, `${req.headers.user_id}-avatar.${extension}`);
}
