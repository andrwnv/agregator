import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { CommentEntity } from '../model/comment.entity';
import { Repository } from 'typeorm';
import { HttpService } from '@nestjs/axios';
import { firstValueFrom } from 'rxjs';


@Injectable()
export class CommentService {
    constructor(@InjectRepository(CommentEntity) private repo: Repository<CommentEntity>,
                private httpService: HttpService) {
    }

    async getUserInfo(bearerToken: string) {
        return (await firstValueFrom(this.httpService.get('http://127.0.0.1:3010/auth/who_am_i',
            {
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${bearerToken}`,
                },
            }))).data;
    }

    async createComment(commentData: string, user: any) {
        const newData = {
            userId: user.id,
            commentContext: commentData,
        };

        return await this.repo.save(this.repo.manager.create(CommentEntity, newData));
    }

    async deleteComment(commentId: string) {
        return await this.repo.delete({
            id: commentId,
        }).then(() => true).catch(() => false);
    }

    async updateComment(commentId: string, newCommentText: string) {
        const comment = await this.repo.find({
            id: commentId,
        });

        return this.repo.save({
            ...comment[0],
            commentContext: newCommentText,
        }).then(res => (res));
    }

    async getAllComments() {
        return await this.repo.find();
    }
}
