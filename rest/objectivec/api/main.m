//
//  main.m
//  api
//
//  Created by Admin on 14.02.16.
//  Copyright © 2016 exmo. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "ExmoApiHandler.h"

ExmoApiHandler *apiHandler;

int main(int argc, const char * argv[]) {
    @autoreleasepool {
        apiHandler = [[ExmoApiHandler alloc] init];
        
        NSMutableDictionary *post = [[NSMutableDictionary alloc] init];
        
        NSLog(@"user_info:\n");
        NSData *response = [apiHandler getResponseFromServerForPost:post method:@"user_info"];
        NSLog(@"%@\n",[[NSString alloc] initWithData:response encoding:NSUTF8StringEncoding]);
        
        NSLog(@"user_cancelled_orders:\n");
        NSMutableDictionary *post2 = [[NSMutableDictionary alloc] init];
        
        [post2 setObject:@"limit" forKey:@"100"];
        [post2 setObject:@"offset" forKey:@"0"];
        
        NSData *response2 = [apiHandler getResponseFromServerForPost:post2 method:@"user_cancelled_orders"];
        NSLog(@"%@\n",[[NSString alloc] initWithData:response2 encoding:NSUTF8StringEncoding]);
        int key;
        // insert code here...
        
        while (1)
        {
            NSLog(@"Press any key (q to quit):");
            fpurge(stdin);
            key = getc(stdin);
            if (key == 'q')
            {
                break;
            }
            
            NSLog(@"\nYou pressed: %c", (char)key);
        }
        
    }
    return 0;
}
