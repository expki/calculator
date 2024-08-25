import { expect } from 'chai';
import { Compress } from './compress.ts';

const testData1 = new Uint8Array([0,1,2,3,4,5,6,7,8,9,10]);
const testData2 = new Uint8Array([10,11,12,13,14,15,16,17,18,19]);

describe('Compress', () => {
  it('should compress without error', () => {
    const [result, err] = Compress(testData2);
    expect(err).to.be.undefined;
    expect(result).to.to.be.instanceOf(Uint8Array);
  });
  it('should compress twice without error', () => {
    const [result1, err1] = Compress(testData1);
    const [result2, err2] = Compress(testData2);
    expect(err1).to.be.undefined;
    expect(result1).to.be.instanceOf(Uint8Array);
    expect(err2).to.be.undefined;
    expect(result2).to.be.instanceOf(Uint8Array);
    expect(result1).to.not.deep.equal(result2);
  });
});
