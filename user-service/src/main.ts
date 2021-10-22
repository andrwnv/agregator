import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { configService } from './config/config.service';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';

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
        console.log("User service listening on 3010");
    });
}

bootstrap();
