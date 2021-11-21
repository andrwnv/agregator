import { Module } from '@nestjs/common';
import { CommentModule } from './comment/comment.module';
import { TypeOrmModule } from '@nestjs/typeorm';

@Module({
    imports: [
        TypeOrmModule.forRoot({
            type: 'postgres',
            host: 'localhost',
            port: 5432,
            username: 'postgres',
            password: '852456',
            database: 'comment_service_repo',
            entities: ['dist/**/*.entity{.ts,.js}'],
            synchronize: true,

            migrationsTableName: 'migration',
            migrations: ['/src/migration/*.ts'],
            cli: {
                migrationsDir: '/src/migration',
            },
        }),
        CommentModule,
    ],
    controllers: [],
    providers: [],
})
export class AppModule {
}
