import { Body, Controller, Post, Get } from '@nestjs/common';
import { Client, ClientProxy, Transport } from '@nestjs/microservices';
import { IUser } from './interfaces/user.interface';

@Controller('user')
export class UserController {
    @Client({
        transport: Transport.TCP,
        options: {
            port: 3010
        }
    })
    client: ClientProxy;

    async onApplicationBootstrap() {
        await this.client.connect();
    }

    @Post('/')
    createUser(@Body() user: IUser) {
        console.log(user);
        return this.client.send('create.user', user);
    }

    @Get('/')
    getAll() {
        return this.client.send('get.user.list', '');
    }
}
