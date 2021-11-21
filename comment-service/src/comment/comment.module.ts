import { Module } from '@nestjs/common';
import { CommentGateway } from './comment.gateway';
import { TypeOrmModule } from '@nestjs/typeorm';
import { CommentEntity } from '../model/comment.entity';
import { HttpModule } from '@nestjs/axios';
import { CommentService } from './comment.service';

@Module({
    providers: [CommentGateway, CommentService],
    exports: [CommentService],
    imports: [
        HttpModule,
        TypeOrmModule.forFeature([CommentEntity])
    ]
})
export class CommentModule {}
