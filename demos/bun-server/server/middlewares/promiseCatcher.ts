// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function promiseCatcher(err: any, req: any, res: any) {    
    // If err has no specified error code, set error code to 'Internal Server Error (500)'
    if (!err.statusCode) {
        err.statusCode = 500;
    }

    console.error(err)
  
    res.status(err.statusCode).json({
        status: err.statusCode,
        error: err.message
    });
  
  }