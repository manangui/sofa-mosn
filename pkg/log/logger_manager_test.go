/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package log

import (
	"context"
	"testing"

	"github.com/alipay/sofa-mosn/pkg/types"
)

func TestUpdateLoggerConfig(t *testing.T) {
	// reset for test
	errorLoggerManagerInstance.managers = make(map[string]ErrorLogger)
	loggers = make(map[string]*Logger)
	//
	logName := "/tmp/mosn/test_update_logger.log"
	if lg, err := GetOrCreateDefaultErrorLogger(logName, DEBUG); err != nil {
		t.Fatal(err)
	} else {
		if lg.(*errorLogger).level != DEBUG {
			t.Fatal("logger created, but level is not expected")
		}
	}
	if lg, err := GetOrCreateDefaultErrorLogger(logName, INFO); err != nil {
		t.Fatal(err)
	} else {
		if lg.(*errorLogger).level != DEBUG {
			t.Fatal("expected get a logger, not create a new one")
		}
	}
	// keeps the logger
	lg, _ := GetOrCreateDefaultErrorLogger(logName, RAW)
	logger := lg.(*errorLogger)
	if ok := UpdateErrorLoggerLevel("not_exists", INFO); ok {
		t.Fatal("update a not exists logger, expected failed")
	}
	// update log level, effects the logger
	if ok := UpdateErrorLoggerLevel(logName, TRACE); !ok {
		t.Fatal("update logger failed")
	} else {
		if logger.level != TRACE {
			t.Fatal("update logger failed")
		}
	}
	// test disable/ enable
	if ok := ToggleLogger(logName, true); !ok {
		t.Fatal("disable logger failed")
	} else {
		if !logger.Logger.disable {
			t.Fatal("disbale logger failed")
		}
	}
	if ok := ToggleLogger(logName, false); !ok {
		t.Fatal("enable logger failed")
	} else {
		if logger.Logger.disable {
			t.Fatal("enable logger failed")
		}
	}
	// Toggle Logger (not error logger)
	baseLoggerPath := "/tmp/mosn/base_logger.log"
	baseLogger, err := GetOrCreateLogger(baseLoggerPath)
	if err != nil || baseLogger.disable {
		t.Fatalf("Create Logger not expected, error: %v, logger state: %v", err, baseLogger.disable)
	}
	if ok := ToggleLogger(baseLoggerPath, true); !ok {
		t.Fatal("enable base logger failed")
	}
	if !baseLogger.disable {
		t.Fatal("disable Logger failed")
	}

}

func TestByContext(t *testing.T) {
	logName := "/tmp/mosn/test_by_context_default.log"
	if err := InitDefaultLogger(logName, DEBUG); err != nil {
		t.Fatal("init default logger failed")
	}
	if lg := ByContext(nil); lg == nil {
		t.Fatal("by context get a nil logger")
	}
	p := "/tmp/mosn/test_by_context.log"
	lg, _ := GetOrCreateDefaultErrorLogger(p, INFO)
	ctx := context.WithValue(context.Background(), types.ContextKeyLogger, lg)
	logger := ByContext(ctx).(*errorLogger)
	if !(logger.output == p && logger.level == INFO) {
		t.Fatal("bycontext get logger is not expected")
	}
}
