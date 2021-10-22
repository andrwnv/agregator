import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { configService } from './config/config.service';
import { UserModule } from './user/user.module';
import { AuthModule } from './auth/auth.module';

@Module({
    imports: [
        TypeOrmModule.forRoot(configService.getTypeOrmConfig()),
        UserModule,
        AuthModule
    ],
})
export class AppModule {
}
