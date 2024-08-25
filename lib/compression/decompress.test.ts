import { expect } from 'chai';
import { Compress } from './compress.ts';
import { Decompress } from './decompress.ts';

const testData1 = new Uint8Array([0,1,2,3,4,5,6,7,8,9,10]);
const testData2 = new Uint8Array([10,11,12,13,14,15,16,17,18,19]);
const [testDataCompressed1, compressError1] = Compress(testData1);
const [testDataCompressed2, compressError2] = Compress(testData2);
const skip = compressError1 !== undefined || compressError2 !== undefined;

describe('Decompress', () => {
  (skip ? it.skip : it)('should decompress without error', () => {
    const [result, err] = Decompress(testDataCompressed2);
    expect(err).to.be.undefined;
    expect(result).to.be.instanceOf(Uint8Array);
    expect(result).to.deep.equal(testData2);
  });
  (skip ? it.skip : it)('should decompress twice without error', () => {
    const [result1, err1] = Decompress(testDataCompressed1);
    const [result2, err2] = Decompress(testDataCompressed2);
    expect(err1).to.be.undefined;
    expect(result1).to.be.instanceOf(Uint8Array);
    expect(result1).to.deep.equal(testData1);
    expect(err2).to.be.undefined;
    expect(result2).to.be.instanceOf(Uint8Array);
    expect(result2).to.deep.equal(testData2);
    expect(result1).to.not.deep.equal(result2);
  });
});
