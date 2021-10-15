import { Controller } from '@nestjs/common';
import { UserService } from './user.service';
import { MessagePattern, Payload } from '@nestjs/microservices';
import { IUser } from './interfaces/user.interface';

@Controller('user')
export class UserController {
    constructor(private userService: UserService) {}

    @MessagePattern('get.user.list')
    getAll() {
        return this.userService.getAll();
    }

    @MessagePattern('create.user')
    create(@Payload() msg: IUser) {
        return this.userService.createUser(msg);
    }
}
