//
//  main.swift
//  api
//
//  Created by Admin on 14.02.16.
//  Copyright © 2016 exmo. All rights reserved.
//

import Foundation

let api = ApiHandler()
let result = api.userInfo()
let dataString = String(data: result!, encoding: NSUTF8StringEncoding)
NSLog(dataString!)

let result2 = api.canceledOrders(100, offset: 0)
let dataString2 = String(data: result2!, encoding: NSUTF8StringEncoding)
NSLog(dataString2!)