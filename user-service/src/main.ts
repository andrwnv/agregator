import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { configService } from './utils/config/config.service';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';
import { Logger } from '@nestjs/common';

async function bootstrap() {
    const app = await NestFactory.create(AppModule);

    if (!configService.isProduction()) {
        const doc = SwaggerModule.createDocument(app, new DocumentBuilder()
            .setTitle('User API')
            .setVersion('1.0')
            .build()
        );

        SwaggerModule.setup('api', app, doc);
    }

    app.enableCors();

    await app.listen(3010, () => {
        new Logger(AppModule.name).log(`User service listening on 3010`);
    });
}

bootstrap();
