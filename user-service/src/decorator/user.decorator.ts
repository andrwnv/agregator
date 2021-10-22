export interface User {
    id: string;
    username: string;
    password: string;
    age: number;
    eMail: string;
    firstName: string;
    lastName: string;
}

import { createParamDecorator } from '@nestjs/common';

// TODO: add auth
export const User = createParamDecorator((data, req) => {
    return req.user || {
        id: 'WAITING PASSPORT JS',
        username: 'WAITING PASSPORT JS',
        password: 'WAITING PASSPORT JS',
        age: null,
        eMail: 'WAITING PASSPORT JS',
        firstName: 'WAITING PASSPORT JS',
        lastName: 'WAITING PASSPORT JS',
    };
});
