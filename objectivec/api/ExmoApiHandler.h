//
//  BtceApiHandler.h
//  
//

#import <Foundation/Foundation.h>

@interface ExmoApiHandler : NSObject

@property (nonatomic,strong) NSString *api_key;
@property (nonatomic,strong) NSString *secret_key;

- (NSData *)getResponseFromServerForPost:(NSDictionary *)postDictionary method:(NSString *)methodName;

- (NSData *)getResponseFromPublicServerUrl:(NSString *)urlString;

@end
