import { Injectable, CanActivate, ExecutionContext } from '@nestjs/common';
import { Reflector } from '@nestjs/core';

import { UserRoles } from './roles.enum';
import { ROLES_KEY } from './roles.decorator';


@Injectable()
export class StaffAccessGuard implements CanActivate {
    constructor(private reflector: Reflector) {
    }

    canActivate(context: ExecutionContext): boolean {
        const requiredRoles = this.reflector.getAllAndOverride<UserRoles[]>(ROLES_KEY, [
            context.getHandler(),
            context.getClass(),
        ]);

        if (!requiredRoles)
            return true;

        // TODO(andrwnv): fix it after creating authentication service.
        const userRole = context.switchToHttp().getRequest().headers['user_role'];

        return requiredRoles.some((role) => role === userRole);
    }
}
