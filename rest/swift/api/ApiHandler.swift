//
//  ApiHandler.swift
//  api
//
//  Created by Admin on 14.02.16.
//  Copyright © 2016 exmo. All rights reserved.
//

import Foundation
import AppKit

public class ApiHandler: NSObject {
    private enum Config: String{
        case API_URL = "https://api.exmo.com/v1/"
        case API_KEY = "your_key"
        case API_SECRET = "your_secret"
        case NONCE = "Nonce"
    }
    
    

    private var api_key: String!
    private var secret_key: String!
        
    
        
    private var nonce: Int{
        get{
            let value = NSUserDefaults.standardUserDefaults().integerForKey(Config.NONCE.rawValue)
            return (value == 0) ? calculateInitialNonce(): value
        }
        set{
            NSUserDefaults.standardUserDefaults().setObject(newValue, forKey: Config.NONCE.rawValue)
        }
    }
        
    override init() {
        super.init()
        setupInitValues()
    }
        
    internal func setupInitValues(){
        self.api_key = Config.API_KEY.rawValue
        self.secret_key = Config.API_SECRET.rawValue
    }
    
    public func userInfo()-> NSData?{
        NSLog("start user_info")
        let post = NSMutableDictionary()
        return self.getResponseFromServerForPost(post, method:"user_info")
    }
    
    public func canceledOrders(limit:Int, offset:Int)-> NSData?{
        NSLog("start user_cancelled_orders")
        let post = NSMutableDictionary()
        post.setObject(limit, forKey: "limit")
        post.setObject(offset, forKey: "offset")
        return self.getResponseFromServerForPost(post, method:"user_cancelled_orders")
    }
    
    private func getResponseFromServerForPost(postDictionary: NSDictionary, method: String) -> NSData?{
            var post: String!
            var i: Int = 0
            for (key, value) in postDictionary {
                if (i==0){
                    post = "\(key)=\(value)"
                }else{
                    post = "\(post)&\(key)=\(value)"
                }
                i++;
            }
            post = "\(post)&nonce=\(nonce)"
            nonce++
            print(post)
            let signedPost = hmacForKeyAndData(secret_key, data: post) as String
            let request = NSMutableURLRequest(URL: NSURL(string: Config.API_URL.rawValue as String + method)!)
            request.HTTPMethod = "POST"
            request.setValue(api_key, forHTTPHeaderField: "Key")
            request.setValue(signedPost, forHTTPHeaderField: "Sign")
            
            let requestBodyData = (post as NSString).dataUsingEncoding(NSUTF8StringEncoding)
            request.HTTPBody = requestBodyData
            
            var error: NSError?
            let theResponse: AutoreleasingUnsafeMutablePointer <NSURLResponse?>=nil
            let responseData = try! NSURLConnection.sendSynchronousRequest(request, returningResponse: theResponse) as NSData!
            if (error != nil){
                return nil
            }
        
            return responseData
        }
        
        private func calculateInitialNonce()->Int{
            let dataFormat = NSDateFormatter()
            dataFormat.dateFormat = "yyyy-MM-dd HH:mm:ss xxxx"
            let timeStamp = NSDate().timeIntervalSinceDate(dataFormat.dateFromString("2012-04-18 00:00:03 +0600")!)
            let currentNonce = NSNumber(double: timeStamp) as Int
            return currentNonce
        }
    
        
        private func hmacForKeyAndData(key: NSString, data: NSString)->NSString{
            let cKey =  key.cStringUsingEncoding(NSASCIIStringEncoding)
            let cData = data.cStringUsingEncoding(NSASCIIStringEncoding)
            let _ = [CUnsignedChar](count: Int(CC_SHA512_DIGEST_LENGTH), repeatedValue: 0)
            let digestLen = Int(CC_SHA512_DIGEST_LENGTH)
            let result = UnsafeMutablePointer<CUnsignedChar>.alloc(digestLen)
            print("CCHmac")
            CCHmac(CCHmacAlgorithm(kCCHmacAlgSHA512), cKey, Int(key.length), cData, Int(data.length), result)
            let hashString =  NSMutableString(capacity: Int(CC_SHA512_DIGEST_LENGTH))
            for i in 0..<digestLen{
                hashString.appendFormat("%02x", result[i])
            }
            return hashString
        }
}
