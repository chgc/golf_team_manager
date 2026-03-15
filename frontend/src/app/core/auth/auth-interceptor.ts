import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';

import { AuthShell } from './auth-shell';

export const authInterceptor: HttpInterceptorFn = (request, next) => {
  const authShell = inject(AuthShell);
  const token = authShell.token();

  if (!token) {
    return next(request);
  }

  return next(
    request.clone({
      setHeaders: {
        Authorization: `Bearer ${token}`,
      },
    }),
  );
};
