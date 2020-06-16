//
//  ExmoApiHandler.m
//
//

#import "ExmoApiHandler.h"
#import <CommonCrypto/CommonHMAC.h>

@implementation ExmoApiHandler

@synthesize api_key;
@synthesize secret_key;

- (id)init {
    self = [super init];
    if (self) {
        [self setupInitialValues];
    }
    return self;
}

- (void)setupInitialValues {
    api_key = @"your_key";
    secret_key = @"your_secret";
}

//

- (NSData *)getResponseFromServerForPost:(NSDictionary *)postDictionary method:(NSString *)methodName{
    NSString *post;
    int i = 0;
    for (NSString *key in [postDictionary allKeys]) {
        NSString *value = [postDictionary objectForKey:key];
        if (i==0)
            post = [NSString stringWithFormat:@"%@=%@", key, value];
        else
            post = [NSString stringWithFormat:@"%@&%@=%@", post, key, value];
        i++;
    }
    post = [NSString stringWithFormat:@"%@&nonce=%@", post, getNonce()];
    
    
    NSString *signedPost = hmacForKeyAndData(secret_key, post);
    NSString *url = [@"https://api.exmo.com/v1/" stringByAppendingString:methodName];
    NSMutableURLRequest *request = [[NSMutableURLRequest alloc]
                                    initWithURL:
                                    [NSURL URLWithString: url]];
    [request setHTTPMethod:@"POST"];
    [request setValue:api_key forHTTPHeaderField:@"Key"];
    [request setValue:signedPost forHTTPHeaderField:@"Sign"];
    [request setHTTPBody:[post dataUsingEncoding: NSUTF8StringEncoding]];
    
    NSURLResponse *theResponse = NULL;
    NSError *theError = NULL;
    NSData *responseData = [NSURLConnection sendSynchronousRequest:request returningResponse:&theResponse error:&theError];
    return responseData;
}

- (NSData *)getResponseFromPublicServerUrl:(NSString *)urlString {
    NSURL *url = [NSURL URLWithString:urlString];
    NSData *data = [NSData dataWithContentsOfURL:url];
    return data;
}

NSString *getNonce() {
    NSTimeInterval timeStamp = [[NSDate date] timeIntervalSinceDate:[NSDate dateWithString:@"2012-04-18 00:00:01 +0600"]];
    int currentNonce = [NSNumber numberWithDouble: timeStamp].intValue;
    NSString *nonceString = [NSString stringWithFormat:@"%i",currentNonce];
    return nonceString;
}

NSString *hmacForKeyAndData(NSString *key, NSString *data) {
    const char *cKey = [key cStringUsingEncoding:NSASCIIStringEncoding];
    const char *cData = [data cStringUsingEncoding:NSASCIIStringEncoding];
    unsigned char cHMAC[CC_SHA512_DIGEST_LENGTH];
    CCHmac(kCCHmacAlgSHA512, cKey, strlen(cKey), cData, strlen(cData), cHMAC);
    NSMutableString *hashString = [NSMutableString stringWithCapacity:sizeof(cHMAC) * 2];
    for (int i = 0; i < sizeof(cHMAC); i++) {
        [hashString appendFormat:@"%02x", cHMAC[i]];
    }
    return hashString;
}

@end
