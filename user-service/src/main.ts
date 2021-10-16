import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';

async function bootstrap() {
    const app = await NestFactory.create(AppModule);
    await app.listen(3010, () => {
        console.log("User service listening on 3010");
    });
}

bootstrap();
