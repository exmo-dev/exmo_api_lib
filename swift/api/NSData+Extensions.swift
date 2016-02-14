//
//  NSData+Extensions.swift
//  api
//
//  Created by Admin on 14.02.16.
//  Copyright © 2016 exmo. All rights reserved.
//
import Foundation

extension NSData{
    func getJsonFromNSData() -> NSDictionary?{
        if let json = try! NSJSONSerialization.JSONObjectWithData(self, options: NSJSONReadingOptions.MutableContainers) as? NSDictionary{
            return json
        }else{
            return nil
        }
    }
}