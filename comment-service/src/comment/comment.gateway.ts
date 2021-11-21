import { MessageBody, SubscribeMessage, WebSocketGateway, WebSocketServer, WsResponse } from '@nestjs/websockets';

import { Server } from 'socket.io';
import { CommentService } from './comment.service';


@WebSocketGateway({
    cors: true
})
export class CommentGateway {
    constructor(private commentService: CommentService) {
    }

    @WebSocketServer()
    server: Server;

    @SubscribeMessage('comment:create')
    async create(@MessageBody() data: any): Promise<WsResponse> {
        const user = await this.commentService.getUserInfo(data.bearerToken);
        const newComment = await this.commentService.createComment(data.commentContext, user);

        return {
            event: 'comment:create_receive',
            data: {
                id: newComment.id,
                comment: newComment.commentContext,
                user: user,
            },
        };
    }

    @SubscribeMessage('comment:delete')
    async delete(@MessageBody() data: any): Promise<WsResponse> {
        return {
            event: 'comment:delete_receive',
            data: {
                success: await this.commentService.deleteComment(data.commentId),
            },
        };
    }

    @SubscribeMessage('comment:update')
    async update(@MessageBody() data: any): Promise<WsResponse> {
        return {
            event: 'comment:update_receive',
            data: await this.commentService.updateComment(data.commentId, data.commentContext),
        };
    }

    @SubscribeMessage('comment:get_all')
    async getAll(): Promise<WsResponse> {
        return {
            event: 'comment:get_all_receive',
            data: await this.commentService.getAllComments(),
        };
    }
}
